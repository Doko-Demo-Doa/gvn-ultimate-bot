// Note: This is a pretty complicated component
import {
  ActionIcon,
  Avatar,
  Button,
  ColorPicker,
  Divider,
  Fieldset,
  Group,
  Image,
  Popover,
  SegmentedControl,
  Space,
  Stack,
  Text,
  Textarea,
  TextInput,
} from "@mantine/core";
import { IconPencil, IconPlus, IconSun } from "@tabler/icons-react";
import { useForm } from "@mantine/form";
import * as classes from "./embed-editor.css";
import { useState } from "react";

interface Props {
  messageId: string;
}

// https://github.com/skyra-project/discord-components
const EmbedEditor: React.FC<Props> = () => {
  const [embedEnabled, setEmbedEnabled] = useState(true);
  const [value, setValue] = useState("emoji");

  const form = useForm({
    mode: "uncontrolled",
    initialValues: {
      color: "#fff",
      mainMessage: "",
      headerMessage: "",
      titleMessage: "",
      embedMainMessage: "",
    },
  });

  return (
    <Stack>
      <Stack>
        <Text>1. Create a message</Text>
        <TextInput
          width={500}
          placeholder="Write your message here!"
          {...form.getInputProps("mainMessage")}
        />
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
                <ColorPicker
                  format="hex"
                  {...form.getInputProps("color")}
                  onChange={(newColor) => form.setFieldValue("color", newColor)}
                />
              </Popover.Dropdown>
            </Popover>

            <ActionIcon color="blue" size="lg" variant="outline">
              <IconPencil size="1rem" />
            </ActionIcon>
          </Stack>

          <Stack>
            <Group
              align="start"
              className={classes.groupWrapper}
              style={{ borderLeftColor: form.getValues().color }}
            >
              <Stack>
                <Group>
                  <Avatar radius="xl" />
                  <TextInput
                    placeholder="Header"
                    {...form.getInputProps("headerMessage")}
                  />
                </Group>
                <Group>
                  <Stack>
                    <TextInput placeholder="Title" />
                    <Textarea placeholder="Main message" />
                  </Stack>

                  <Image
                    radius="md"
                    h={100}
                    w="auto"
                    fit="contain"
                    src="https://raw.githubusercontent.com/mantinedev/mantine/master/.demo/images/bg-9.png"
                  />
                </Group>

                {[1].map((i, n) => (
                  <Fieldset key={i} legend={`Custom field ${i}`} disabled>
                    <TextInput placeholder="Field name" />
                    <TextInput placeholder="Field value" mt="md" />
                  </Fieldset>
                ))}
                <Button leftSection={<IconPlus size={14} />} variant="default">
                  Add new field
                </Button>

                <Divider variant="dotted" my="md" />

                <Group>
                  <Avatar radius="xl" />
                  <TextInput placeholder="Footer" />
                </Group>
              </Stack>
            </Group>
          </Stack>
        </Group>
      )}
    </Stack>
  );
};

export default EmbedEditor;
