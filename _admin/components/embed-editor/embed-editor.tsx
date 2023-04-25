// Note: This is a pretty complicated component

import { Box, Button, ColorPicker, Popover, Text } from "@mantine/core";

interface Props {
  messageId: string;
}

const EmbedEditor: React.FC<Props> = () => {
  return (
    <Box>
      <Popover closeOnClickOutside position="bottom" withArrow shadow="md">
        <Popover.Target>
          <Button compact>Toggle popover</Button>
        </Popover.Target>
        <Popover.Dropdown>
          <ColorPicker format="hex" />
        </Popover.Dropdown>
      </Popover>
    </Box>
  );
};

export default EmbedEditor;
