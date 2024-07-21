"use client";

// Note: This is a pretty complicated component
import {
  ActionIcon,
  Box,
  ColorPicker,
  Group,
  Paper,
  Popover,
  SegmentedControl,
  Space,
  Stack,
  Text,
  TextInput,
} from "@mantine/core";
import { IconPencil, IconSun } from "@tabler/icons-react";
import * as classes from "./embed-editor.css";
import { useState } from "react";

interface Props {
  messageId: string;
}

const EmbedEditor: React.FC<Props> = () => {
  const [embedEnabled, setEmbedEnabled] = useState(false);
  const [value, setValue] = useState("emoji");

  return (
    <Stack>
      <Stack>
        <Text>1. Create a message</Text>
        <TextInput width={300} placeholder="Write your message here!" />
      </Stack>

      <Space h="lg" />

      <Stack>
        <Text>2. Add reactions</Text>

        <SegmentedControl
          value={value}
          onChange={setValue}
          className={classes.segmented}
          data={[
            { label: "Emoji", value: "emoji" },
            { label: "Button", value: "button" },
            { label: "Dropdown", value: "dropdown" },
          ]}
        />
      </Stack>

      {embedEnabled && (
        <Group align="start">
          <Stack>
            <Popover
              closeOnClickOutside
              position="bottom"
              withArrow
              shadow="md"
            >
              <Popover.Target>
                <ActionIcon color="red.6" size="lg" variant="outline">
                  <IconSun size="1rem" />
                </ActionIcon>
              </Popover.Target>

              <Popover.Dropdown>
                <ColorPicker format="hex" />
              </Popover.Dropdown>
            </Popover>

            <ActionIcon color="blue" size="lg" variant="outline">
              <IconPencil size="1rem" />
            </ActionIcon>
          </Stack>

          <Stack>
            <Paper>
              <Box className={classes.typeIndicator} />
            </Paper>
          </Stack>
        </Group>
      )}
    </Stack>
  );
};

export default EmbedEditor;
