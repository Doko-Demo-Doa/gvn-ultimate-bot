// Note: This is a pretty complicated component
import {
  ActionIcon,
  AspectRatio,
  Avatar,
  Box,
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
  TextInput,
  Textarea,
} from "@mantine/core";
import { useForm } from "@mantine/form";
import { notifications } from "@mantine/notifications";
import { IconPencil, IconPlus, IconSun } from "@tabler/icons-react";
import { useState } from "react";
import { v4 as uuidv4 } from "uuid";
import * as classes from "~/components/embed-editor/embed-editor.css";
import { UploadDropzone } from "~/utils/uploadthing";

interface Props {
  messageId: string;
}

interface IFormData {
  color: string;
  mainMessage: string;
  headerMessage: string;
  titleMessage: string;
  embedMainMessage: string;
  featuredImage: string;
  customFields: Array<{ id: string; fieldName: string; fieldValue: string }>;
}

const MAX_CUSTOM_FIELDS = 5;

// https://github.com/skyra-project/discord-components
const EmbedEditor: React.FC<Props> = () => {
  const [value, setValue] = useState("emoji");
  const [files, setFiles] = useState<File[]>([]);
  const [mainImageUrl, setMainImageUrl] = useState("");

  const form = useForm<IFormData>({
    mode: "uncontrolled",
    initialValues: {
      color: "#fff",
      mainMessage: "",
      headerMessage: "",
      titleMessage: "",
      embedMainMessage: "",
      featuredImage: "",
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
      <Stack className="embed-editor-wrapper">
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

        <Group align="stretch" className="main-area">
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
          </Stack>

          <Stack className={classes.rightArea}>
            <Group
              align="start"
              className={classes.groupWrapper}
              style={{ borderLeftColor: form.getValues().color }}
            >
              <Stack className={classes.mainEmbedEditorArea}>
                <Group>
                  <Avatar radius="xl" />
                  <TextInput
                    placeholder="Header"
                    {...form.getInputProps("headerMessage")}
                  />
                </Group>
                <Group className={classes.titleAndMainContentArea}>
                  <Stack className={classes.titleAndMainText}>
                    <TextInput
                      placeholder="Title"
                      {...form.getInputProps("titleMessage")}
                    />
                    <Textarea
                      placeholder="Main content"
                      maxRows={9}
                      autosize
                      {...form.getInputProps("embedMainMessage")}
                    />
                  </Stack>

                  <AspectRatio className={classes.mainImageWrapper}>
                    {mainImageUrl && (
                      <Image
                        h={250}
                        className={classes.mainImage}
                        w="100%"
                        src={mainImageUrl}
                      />
                    )}
                    <Box className={classes.uploadBoxWrapper}>
                      <UploadDropzone
                        endpoint="imageUploader"
                        config={{ mode: "auto" }}
                        content={{ allowedContent: "image/png,image/jpeg" }}
                        onBeforeUploadBegin={(files) => {
                          // Preprocess files before uploading (e.g. rename them)
                          return files.map(
                            (f) =>
                              new File([f], `temp-${f.name}`, {
                                type: f.type,
                              })
                          );
                        }}
                        onClientUploadComplete={(res) => {
                          // Do something with the response
                          setMainImageUrl(res[0].url);
                        }}
                        onUploadError={(error: Error) => {
                          // Do something with the error.
                          console.log("Error Uploadthing: ", error);
                          notifications.show({
                            color: "red",
                            title: "Lỗi",
                            message: "Không thể upload file này",
                          });
                        }}
                      />
                    </Box>
                  </AspectRatio>
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
      </Stack>
    </form>
  );
};

export default EmbedEditor;
