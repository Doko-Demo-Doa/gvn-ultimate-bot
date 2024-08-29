// Note: This is a pretty complicated component
import {
  ActionIcon,
  Avatar,
  Button,
  ColorPicker,
  Divider,
  Fieldset,
  FileButton,
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
import { v4 as uuidv4 } from "uuid";

interface Props {
  messageId: string;
}

interface IFormData {
  color: string;
  mainMessage: string;
  headerMessage: string;
  titleMessage: string;
  embedMainMessage: string;
  customFields: Array<{ id: string; fieldName: string; fieldValue: string }>;
}

const MAX_CUSTOM_FIELDS = 5;

// https://github.com/skyra-project/discord-components
const EmbedEditor: React.FC<Props> = () => {
  const [embedEnabled, setEmbedEnabled] = useState(true);
  const [value, setValue] = useState("emoji");
  const [files, setFiles] = useState<File[]>([]);

  const form = useForm<IFormData>({
    mode: "uncontrolled",
    initialValues: {
      color: "#fff",
      mainMessage: "",
      headerMessage: "",
      titleMessage: "",
      embedMainMessage: "",
      customFields: [
        {
          id: uuidv4(),
          fieldName: "",
          fieldValue: "",
        },
      ],
    },
  });

  const handleSubmit = (values: typeof form.values) => {
    console.log(values);
  };

  const cFields = form.getValues().customFields;

  return (
    <form onSubmit={form.onSubmit((values) => handleSubmit(values))}>
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
                    onChange={(newColor) =>
                      form.setFieldValue("color", newColor)
                    }
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
                      <TextInput
                        placeholder="Title"
                        {...form.getInputProps("titleMessage")}
                      />
                      <Textarea
                        placeholder="Main content"
                        {...form.getInputProps("embedMainMessage")}
                      />
                    </Stack>

                    <Image
                      radius="md"
                      h={100}
                      w="auto"
                      fit="contain"
                      src="https://raw.githubusercontent.com/mantinedev/mantine/master/.demo/images/bg-9.png"
                    />
                  </Group>

                  {cFields.map((n, i) => (
                    <Fieldset key={n.id} legend={`Custom field ${i + 1}`}>
                      <TextInput placeholder="Field name" />
                      <TextInput placeholder="Field value" mt="md" />
                    </Fieldset>
                  ))}
                  <Button
                    leftSection={<IconPlus size={14} />}
                    variant="default"
                    disabled={cFields.length >= MAX_CUSTOM_FIELDS}
                    onClick={() => {
                      if (form.values.customFields.length >= MAX_CUSTOM_FIELDS)
                        return;

                      form.insertListItem("customFields", {
                        id: uuidv4(),
                        fieldName: "",
                        fieldValue: "",
                      });
                    }}
                  >
                    Add new field
                  </Button>

                  <Divider variant="dotted" my="md" />

                  <Group>
                    <FileButton
                      onChange={(f) => {
                        if (f) {
                          setFiles([f]);
                        }
                      }}
                      accept="image/png,image/jpeg"
                    >
                      {(props) => (
                        <Avatar
                          styles={{ root: { cursor: "pointer" } }}
                          radius="xl"
                          {...props}
                        />
                      )}
                    </FileButton>

                    <TextInput placeholder="Footer" />
                  </Group>
                </Stack>
              </Group>

              <Button variant="gradient" type="submit">
                Save
              </Button>
            </Stack>
          </Group>
        )}
      </Stack>
    </form>
  );
};

export default EmbedEditor;
