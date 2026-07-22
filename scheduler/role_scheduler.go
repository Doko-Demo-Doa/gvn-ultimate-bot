package scheduler

import (
	"container/heap"
	"doko/gvn-ultimate-bot/models"
	"doko/gvn-ultimate-bot/services/discordservice"
	"log"
	"sync"
	"time"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/state"
)

// heapItem represents a single timed role assignment in the priority queue.
type heapItem struct {
	assignmentID uint
	userNativeID string
	roleNativeID string
	expiresAt    time.Time
	index        int // required by container/heap for in-place updates
}

// assignmentHeap implements container/heap.Interface, ordered by expiresAt (earliest first).
type assignmentHeap []*heapItem

func (h assignmentHeap) Len() int           { return len(h) }
func (h assignmentHeap) Less(i, j int) bool { return h[i].expiresAt.Before(h[j].expiresAt) }
func (h assignmentHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
	h[i].index = i
	h[j].index = j
}

func (h *assignmentHeap) Push(x interface{}) {
	item := x.(*heapItem)
	item.index = len(*h)
	*h = append(*h, item)
}

func (h *assignmentHeap) Pop() interface{} {
	old := *h
	n := len(old)
	item := old[n-1]
	old[n-1] = nil // avoid memory leak
	item.index = -1
	*h = old[0 : n-1]
	return item
}

// RoleScheduler manages timed Discord role assignments using a min-heap + single timer.
type RoleScheduler struct {
	mu       sync.Mutex
	heap     assignmentHeap
	timer    *time.Timer
	state    *state.State
	service  discordservice.DiscordService
	guildID  discord.GuildID
}

// NewRoleScheduler creates a scheduler (Start must be called before use).
func NewRoleScheduler(s *state.State, ds discordservice.DiscordService, guildID discord.GuildID) *RoleScheduler {
	return &RoleScheduler{
		state:   s,
		service: ds,
		guildID: guildID,
		heap:    make(assignmentHeap, 0),
	}
}

// Start rebuilds the heap from the database and arms the first timer.
// Must be called once, after the Discord session is ready.
func (rs *RoleScheduler) Start() {
	rs.mu.Lock()
	defer rs.mu.Unlock()

	// 1. Process backlog: revoke anything that expired while offline
	backlog, err := rs.service.GetExpiredRoleAssignments()
	if err != nil {
		log.Println("[scheduler] error fetching backlog:", err)
	} else {
		for _, a := range backlog {
			rs.revokeOne(a)
		}
	}

	// 2. Rebuild heap from active assignments
	active, err := rs.service.GetAllActiveAssignments()
	if err != nil {
		log.Println("[scheduler] error loading active assignments:", err)
		return
	}

	for _, a := range active {
		heap.Push(&rs.heap, &heapItem{
			assignmentID: a.ID,
			userNativeID: a.UserNativeID,
			roleNativeID: a.RoleNativeID,
			expiresAt:    a.ExpirationDate,
		})
	}
	heap.Init(&rs.heap)

	rs.setNextTimerUnsafe()
	log.Printf("[scheduler] started with %d active assignment(s)", len(active))
}

// GrantRole assigns a Discord role to a user for the given duration.
// It immediately calls the Discord API and persists the assignment in the DB.
func (rs *RoleScheduler) GrantRole(userNativeID string, roleNativeID string, duration time.Duration) error {
	// Add role on Discord immediately
	userID := parseSnowflake(userNativeID)
	roleID := parseSnowflakeRole(roleNativeID)
	if userID == 0 || roleID == 0 {
		return nil // invalid IDs, nothing to do
	}
	if err := rs.state.AddRole(rs.guildID, userID, roleID, api.AddRoleData{}); err != nil {
		return err
	}

	// Persist in database
	assignment, err := rs.service.AssignRoleToUser(userNativeID, roleNativeID, duration)
	if err != nil {
		log.Println("[scheduler] failed to persist assignment:", err)
		return err
	}

	// Push into heap and maybe reset timer
	rs.mu.Lock()
	defer rs.mu.Unlock()

	heap.Push(&rs.heap, &heapItem{
		assignmentID: assignment.ID,
		userNativeID: assignment.UserNativeID,
		roleNativeID: assignment.RoleNativeID,
		expiresAt:    assignment.ExpirationDate,
	})

	if rs.heap[0].assignmentID == assignment.ID {
		// This new assignment is now the soonest-expiring one
		rs.setNextTimerUnsafe()
	}

	return nil
}

// AddRole assigns a Discord role to a user permanently (no DB tracking).
// Used by reaction-role modules where persistence is handled by Discord itself.
func (rs *RoleScheduler) AddRole(userNativeID string, roleNativeID string) error {
	userID := parseSnowflake(userNativeID)
	roleID := parseSnowflakeRole(roleNativeID)
	if userID == 0 || roleID == 0 {
		return nil
	}
	return rs.state.AddRole(rs.guildID, userID, roleID, api.AddRoleData{})
}

// RemoveRole removes a Discord role from a user permanently (no DB tracking).
func (rs *RoleScheduler) RemoveRole(userNativeID string, roleNativeID string) error {
	userID := parseSnowflake(userNativeID)
	roleID := parseSnowflakeRole(roleNativeID)
	if userID == 0 || roleID == 0 {
		return nil
	}
	return rs.state.RemoveRole(rs.guildID, userID, roleID, api.AuditLogReason(""))
}

// RevokeRole performs an early/manual revocation of a timed role assignment.
func (rs *RoleScheduler) RevokeRole(assignmentID uint) error {
	// Fetch the assignment to know who/what to revoke
	assignment, err := rs.service.GetAssignmentByID(assignmentID)
	if err != nil {
		return err
	}

	// Remove from Discord
	userID := parseSnowflake(assignment.UserNativeID)
	roleID := parseSnowflakeRole(assignment.RoleNativeID)
	if userID != 0 && roleID != 0 {
		if err := rs.state.RemoveRole(rs.guildID, userID, roleID, api.AuditLogReason("")); err != nil {
			log.Printf("[scheduler] failed to remove role %s from user %s: %v", assignment.RoleNativeID, assignment.UserNativeID, err)
			// Continue: still delete DB record
		}
	}

	// Soft-delete from DB
	if err := rs.service.RevokeRoleAssignment(assignmentID); err != nil {
		return err
	}

	// Rebuild heap from scratch so the cancelled item disappears
	rs.mu.Lock()
	defer rs.mu.Unlock()
	rs.rebuildHeapUnsafe()

	return nil
}

// onTimerFire is called when the timer reaches the next expiry.
func (rs *RoleScheduler) onTimerFire() {
	rs.mu.Lock()
	defer rs.mu.Unlock()

	now := time.Now()

	// Process all expired items (there may be multiple with the same expiry)
	for rs.heap.Len() > 0 {
		top := rs.heap[0]
		if top.expiresAt.After(now) {
			break
		}

		heap.Pop(&rs.heap)

		assignment, err := rs.service.GetAssignmentByID(top.assignmentID)
		if err != nil || assignment == nil {
			// Already deleted (e.g. manual revoke), skip
			continue
		}

		rs.revokeOne(assignment)
	}

	rs.setNextTimerUnsafe()
}

// revokeOne removes a single Discord role and soft-deletes the DB record.
// Called while holding the lock.
func (rs *RoleScheduler) revokeOne(assignment *models.DiscordUserRole) {
	userID := parseSnowflake(assignment.UserNativeID)
	roleID := parseSnowflakeRole(assignment.RoleNativeID)

	if userID == 0 || roleID == 0 {
		log.Printf("[scheduler] invalid IDs for assignment %d, cleaning DB only", assignment.ID)
		_ = rs.service.RevokeRoleAssignment(assignment.ID)
		return
	}

	if err := rs.state.RemoveRole(rs.guildID, userID, roleID, api.AuditLogReason("")); err != nil {
		log.Printf("[scheduler] failed to remove role %s from user %s: %v", assignment.RoleNativeID, assignment.UserNativeID, err)
		// Keep DB record for retry; the item is already popped from the heap,
		// so it won't fire again. On next startup, the backlog scan will retry.
		return
	}

	if err := rs.service.RevokeRoleAssignment(assignment.ID); err != nil {
		log.Printf("[scheduler] failed to revoke assignment %d in DB: %v", assignment.ID, err)
		return
	}

	log.Printf("[scheduler] revoked role %s from user %s (assignment %d)", assignment.RoleNativeID, assignment.UserNativeID, assignment.ID)
}

// setNextTimerUnsafe arms the timer for the heap-top item. Must hold the lock.
func (rs *RoleScheduler) setNextTimerUnsafe() {
	if rs.timer != nil {
		rs.timer.Stop()
	}

	if rs.heap.Len() == 0 {
		return
	}

	top := rs.heap[0]
	delay := time.Until(top.expiresAt)
	if delay < 0 {
		delay = 0
	}

	rs.timer = time.AfterFunc(delay, func() { rs.onTimerFire() })
}

// rebuildHeapUnsafe rebuilds the heap from active DB assignments. Must hold the lock.
func (rs *RoleScheduler) rebuildHeapUnsafe() {
	if rs.timer != nil {
		rs.timer.Stop()
	}
	rs.heap = rs.heap[:0]

	active, err := rs.service.GetAllActiveAssignments()
	if err != nil {
		log.Println("[scheduler] error rebuilding heap:", err)
		rs.setNextTimerUnsafe()
		return
	}

	for _, a := range active {
		rs.heap = append(rs.heap, &heapItem{
			assignmentID: a.ID,
			userNativeID: a.UserNativeID,
			roleNativeID: a.RoleNativeID,
			expiresAt:    a.ExpirationDate,
		})
	}
	heap.Init(&rs.heap)
	rs.setNextTimerUnsafe()
}

func parseSnowflake(val string) discord.UserID {
	s, err := discord.ParseSnowflake(val)
	if err != nil {
		return 0
	}
	return discord.UserID(s)
}

func parseSnowflakeRole(val string) discord.RoleID {
	s, err := discord.ParseSnowflake(val)
	if err != nil {
		return 0
	}
	return discord.RoleID(s)
}
