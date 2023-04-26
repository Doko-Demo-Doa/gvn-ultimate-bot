// Note: This is a pretty complicated component

import {
  ActionIcon,
  Box,
  Button,
  ColorPicker,
  Group,
  Paper,
  Popover,
  Stack,
  Text,
  TextInput,
  Tooltip,
  createStyles,
} from "@mantine/core";
import { IconSun } from "@tabler/icons-react";

interface Props {
  messageId: string;
}

const useStyles = createStyles((theme) => ({
  wrapperFrame: {
    width: "400px",
    height: "400px",
    backgroundColor: "white",
    display: "flex",
    justifyContent: "flex-end",
    borderRadius: 4,
  },
  mainFrame: {
    backgroundColor: "red",
    width: "calc(100% - 4px)",
    height: "100%",
  },
}));

const EmbedEditor: React.FC<Props> = () => {
  const { classes } = useStyles();

  return (
    <Group>
      <Stack sx={{ maxWidth: "2rem" }}>
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

      <Stack sx={{ flexGrow: 2 }}>
        <TextInput placeholder="Write your message here!" />

        <Paper>
          <Box className={classes.wrapperFrame}>
            <Box className={classes.mainFrame} />
          </Box>
        </Paper>
      </Stack>
    </Group>
  );
};

export default EmbedEditor;
