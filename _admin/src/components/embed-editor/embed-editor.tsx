// Note: This is a pretty complicated component
"use client";

import {
  ActionIcon,
  Box,
  ColorPicker,
  Group,
  Paper,
  Popover,
  Stack,
  TextInput,
} from "@mantine/core";
import { IconSun } from "@tabler/icons-react";

interface Props {
  messageId: string;
}

const EmbedEditor: React.FC<Props> = () => {
  return (
    <Group>
      <Stack>
        <Popover closeOnClickOutside position="bottom" withArrow shadow="md">
          <Popover.Target>
            <ActionIcon color="red.6" variant="outline">
              <IconSun size="1rem" />
            </ActionIcon>
          </Popover.Target>
          <Popover.Dropdown>
            <ColorPicker format="hex" />
          </Popover.Dropdown>
        </Popover>
      </Stack>

      <Stack>
        <TextInput placeholder="Write your message here!" />

        <Paper>
          <Box>
            <Box />
          </Box>
        </Paper>
      </Stack>
    </Group>
  );
};

export default EmbedEditor;
