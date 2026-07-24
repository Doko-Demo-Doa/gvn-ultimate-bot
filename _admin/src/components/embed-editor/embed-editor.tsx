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
  Flex,
  Group,
  Image,
  NativeSelect,
  Popover,
  SegmentedControl,
  Select,
  Space,
  Stack,
  Text,
  Textarea,
  TextInput,
} from "@mantine/core";
import { schemaResolver, useForm } from "@mantine/form";
import { notifications } from "@mantine/notifications";
import { IconPlus, IconSun, IconTrashFilled } from "@tabler/icons-react";
import { useState } from "react";
import { v4 as uuidv4 } from "uuid";
import { z } from "zod/v4";
import * as classes from "~/components/embed-editor/embed-editor.css";
import { optimizeImage } from "~/datasource/rest/image-optimizer";
import { vars } from "~/theme";
import type {
  IDiscordChannel,
  IDiscordEmoji,
  IDiscordRole,
  IReactionRoleMessagePayload,
} from "~/types/types";
import { UploadDropzone } from "~/utils/uploadthing";

interface Props {
  roles: IDiscordRole[];
  channels: IDiscordChannel[];
  emojis: IDiscordEmoji[];
  onPublish: (payload: IReactionRoleMessagePayload) => void;
  isPublishing: boolean;
  initialPayload?: IReactionRoleMessagePayload;
}

const MAX_CUSTOM_FIELDS = 5;
const MAX_INTERACTIONS = 10;

const dropdownOptionSchema = z.object({
  id: z.string(),
  label: z.string().min(1, "Label is required"),
  emoji: z.string().optional(),
  description: z.string().optional(),
  role_native_id: z.string().min(1, "Role is required"),
});

const interactionSchema = z.object({
  id: z.string(),
  type: z.enum(["emoji", "button", "dropdown"]),
  emoji: z.string().optional(),
  label: z.string().optional(),
  style: z.enum(["primary", "secondary", "success", "danger"]).optional(),
  role_native_id: z.string().optional(),
  placeholder: z.string().optional(),
  options: z.array(dropdownOptionSchema).optional(),
});

const schema = z.object({
  channel_id: z.string().min(1, "Channel ID is required"),
  mode: z.enum(["default", "reverse"]).default("default"),
  color: z.string().optional(),
  mainMessage: z.string().optional(),
  headerMessage: z.string().optional(),
  footerMessage: z.string().optional(),
  titleMessage: z.string().optional(),
  embedMainMessage: z.string().optional(),
  featuredImage: z.string().optional(),
  customFields: z
    .array(
      z.object({
        id: z.string(),
        fieldName: z.string(),
        fieldValue: z.string(),
      }),
    )
    .max(MAX_CUSTOM_FIELDS),
  interactions: z.array(interactionSchema).max(MAX_INTERACTIONS),
});

type IFormData = z.infer<typeof schema>;

function payloadToFormData(payload?: IReactionRoleMessagePayload): IFormData {
  if (!payload) {
    return {
      channel_id: "",
      mode: "default",
      color: "#5865F2",
      mainMessage: "",
      headerMessage: "",
      footerMessage: "",
      titleMessage: "",
      embedMainMessage: "",
      featuredImage: "",
      customFields: [{ id: uuidv4(), fieldName: "", fieldValue: "" }],
      interactions: [],
    };
  }
  const colorHex = payload.embed?.color
    ? `#${payload.embed.color.toString(16).padStart(6, "0")}`
    : "#5865F2";
  return {
    channel_id: payload.channel_id || "",
    mode: payload.mode || "default",
    color: colorHex,
    mainMessage: payload.message || "",
    headerMessage: payload.embed?.author || "",
    footerMessage: payload.embed?.footer || "",
    titleMessage: payload.embed?.title || "",
    embedMainMessage: payload.embed?.description || "",
    featuredImage: payload.embed?.image_url || "",
    customFields:
      payload.embed?.fields && payload.embed.fields.length > 0
        ? payload.embed.fields.map((f) => ({
            id: uuidv4(),
            fieldName: f.name,
            fieldValue: f.value,
          }))
        : [{ id: uuidv4(), fieldName: "", fieldValue: "" }],
    interactions: (payload.interactions || []).map((it) => ({
      id: it.id,
      type: it.type,
      emoji: it.emoji || undefined,
      label: it.label || undefined,
      style: it.style || undefined,
      role_native_id: it.role_native_id || undefined,
      placeholder: it.placeholder || undefined,
      options: it.options
        ? it.options.map((opt) => ({
            id: opt.id,
            label: opt.label,
            emoji: opt.emoji || undefined,
            description: opt.description || undefined,
            role_native_id: opt.role_native_id,
          }))
        : undefined,
    })),
  };
}

// https://github.com/skyra-project/discord-components
const emojiSelectData = (emojis: IDiscordEmoji[]) => [
  { value: "", label: "-- Không có emoji --", disabled: true },
  ...emojis.map((e) => ({
    value: e.api_name,
    label: `:${e.name}:`,
    image_url: e.image_url,
  })),
];

function renderEmojiOption({ option }: { option: any }) {
  return (
    <Group gap="xs">
      {option.image_url && (
        // biome-ignore lint/performance/noImgElement: small Discord CDN emoji thumbnails in a dropdown; next/Image adds complexity with no benefit here
        <img
          src={option.image_url}
          alt={option.label}
          width={20}
          height={20}
          style={{ objectFit: "contain", display: "inline-block" }}
        />
      )}
      <span>{option.label}</span>
    </Group>
  );
}

const EmbedEditor: React.FC<Props> = ({
  roles,
  channels,
  emojis,
  onPublish,
  isPublishing,
  initialPayload,
}) => {
  const [value, setValue] = useState("emoji");
  const [mainImageUrl, setMainImageUrl] = useState(
    initialPayload?.embed?.image_url || "",
  );

  const form = useForm<IFormData>({
    mode: "uncontrolled",
    initialValues: payloadToFormData(initialPayload),
    validate: schemaResolver(schema, { sync: true }),
  });

  const handleSubmit = (values: typeof form.values) => {
    const payload: IReactionRoleMessagePayload = {
      channel_id: values.channel_id,
      message: values.mainMessage || undefined,
      mode: values.mode,
      embed:
        values.titleMessage ||
        values.embedMainMessage ||
        values.headerMessage ||
        values.footerMessage ||
        values.featuredImage ||
        values.customFields.length > 0
          ? {
              title: values.titleMessage || undefined,
              description: values.embedMainMessage || undefined,
              color:
                parseInt((values.color || "#5865F2").replace("#", ""), 16) ||
                0x5865f2,
              image_url: values.featuredImage || undefined,
              footer: values.footerMessage || undefined,
              author: values.headerMessage || undefined,
              fields: values.customFields
                .filter((f) => f.fieldName && f.fieldValue)
                .map((f) => ({
                  name: f.fieldName,
                  value: f.fieldValue,
                  inline: false,
                })),
            }
          : undefined,
      interactions: values.interactions.map((it) => ({
        id: it.id,
        type: it.type,
        emoji: it.emoji || undefined,
        label: it.label || undefined,
        style: it.style || undefined,
        role_native_id: it.role_native_id || undefined,
        placeholder: it.placeholder || undefined,
        options: it.options || undefined,
      })),
    };

    onPublish(payload);
  };

  const cFields = form.getValues().customFields;
  const interactions = form.getValues().interactions;

  const addInteraction = (type_: "emoji" | "button" | "dropdown") => {
    if (interactions.length >= MAX_INTERACTIONS) return;
    const base = {
      id: uuidv4(),
      type: type_,
    } as const;
    if (type_ === "emoji" || type_ === "button") {
      form.insertListItem("interactions", {
        ...base,
        emoji: "",
        label: "",
        style: "secondary",
        role_native_id: "",
      });
    } else {
      form.insertListItem("interactions", {
        ...base,
        placeholder: "Chọn một role...",
        options: [
          {
            id: uuidv4(),
            label: "",
            role_native_id: "",
          },
        ],
      });
    }
  };

  const removeInteraction = (index: number) => {
    form.removeListItem("interactions", index);
  };

  const roleSelectData = [
    { value: "", label: "-- Chọn role --", disabled: true },
    ...roles.map((r) => ({
      value: r.NativeId,
      label: r.Name,
    })),
  ];

  return (
    <form onSubmit={form.onSubmit((values) => handleSubmit(values))}>
      <Stack className="embed-editor-wrapper">
        <Group>
          <Select
            label="Discord channel"
            placeholder="Chọn channel..."
            searchable
            data={channels.map((ch) => ({
              value: ch.id,
              label: `#${ch.name}`,
            }))}
            {...form.getInputProps("channel_id")}
          />
          <NativeSelect
            label="Chế độ"
            data={[
              { value: "default", label: "Default (cấp role khi tương tác)" },
              { value: "reverse", label: "Reverse (xóa role khi tương tác)" },
            ]}
            {...form.getInputProps("mode")}
          />
        </Group>

        <Divider variant="dotted" my="md" />

        <Stack>
          <Text>1. Message content (optional)</Text>
          <TextInput
            placeholder="Write your message here! This will appear above the embed."
            {...form.getInputProps("mainMessage")}
          />
        </Stack>

        <Space h="lg" />

        <Stack>
          <Text>2. Add interactions</Text>

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

          <Group>
            <Button
              size="xs"
              variant="light"
              onClick={() => addInteraction(value as any)}
              disabled={interactions.length >= MAX_INTERACTIONS}
              leftSection={<IconPlus size={14} />}
            >
              Thêm {value}
            </Button>
          </Group>

          {interactions.map((it, i) => (
            <Fieldset key={it.id} legend={`${it.type.toUpperCase()} ${i + 1}`}>
              <Stack>
                {(it.type === "emoji" || it.type === "button") && (
                  <>
                    <Group>
                      <Select
                        label="Emoji"
                        placeholder="Chọn emoji..."
                        searchable
                        style={{ width: 200 }}
                        data={emojiSelectData(emojis)}
                        renderOption={renderEmojiOption}
                        {...form.getInputProps(`interactions.${i}.emoji`)}
                      />
                      {it.type === "button" && (
                        <TextInput
                          label="Label"
                          placeholder="Nút bấm"
                          style={{ flexGrow: 1 }}
                          {...form.getInputProps(`interactions.${i}.label`)}
                        />
                      )}
                    </Group>
                    {it.type === "button" && (
                      <NativeSelect
                        label="Style"
                        data={[
                          { value: "primary", label: "Primary (blurple)" },
                          { value: "secondary", label: "Secondary (grey)" },
                          { value: "success", label: "Success (green)" },
                          { value: "danger", label: "Danger (red)" },
                        ]}
                        {...form.getInputProps(`interactions.${i}.style`)}
                      />
                    )}
                    <NativeSelect
                      label="Role"
                      data={roleSelectData}
                      {...form.getInputProps(
                        `interactions.${i}.role_native_id`,
                      )}
                    />
                  </>
                )}

                {it.type === "dropdown" && (
                  <>
                    <TextInput
                      label="Placeholder"
                      placeholder="Chọn một role..."
                      {...form.getInputProps(`interactions.${i}.placeholder`)}
                    />
                    <Text size="sm" fw={500}>
                      Options
                    </Text>
                    {it.options?.map((opt, optIndex) => (
                      <Group key={opt.id} align="flex-start">
                        <Stack gap={4} style={{ flexGrow: 1 }}>
                          <Group>
                            <TextInput
                              placeholder="Label"
                              style={{ flexGrow: 1 }}
                              {...form.getInputProps(
                                `interactions.${i}.options.${optIndex}.label`,
                              )}
                            />
                            <Select
                              placeholder="Emoji"
                              searchable
                              style={{ width: 160 }}
                              data={emojiSelectData(emojis)}
                              renderOption={renderEmojiOption}
                              {...form.getInputProps(
                                `interactions.${i}.options.${optIndex}.emoji`,
                              )}
                            />
                          </Group>
                          <TextInput
                            placeholder="Description"
                            {...form.getInputProps(
                              `interactions.${i}.options.${optIndex}.description`,
                            )}
                          />
                          <NativeSelect
                            data={roleSelectData}
                            {...form.getInputProps(
                              `interactions.${i}.options.${optIndex}.role_native_id`,
                            )}
                          />
                        </Stack>
                        <Button
                          size="xs"
                          color="red"
                          variant="subtle"
                          onClick={() => {
                            const current =
                              form.getValues().interactions[i].options || [];
                            const next = current.filter(
                              (_, idx) => idx !== optIndex,
                            );
                            form.setFieldValue(
                              `interactions.${i}.options`,
                              next,
                            );
                          }}
                        >
                          <IconTrashFilled size={14} />
                        </Button>
                      </Group>
                    ))}
                    <Button
                      size="xs"
                      variant="default"
                      leftSection={<IconPlus size={14} />}
                      onClick={() => {
                        const current =
                          form.getValues().interactions[i].options || [];
                        form.setFieldValue(`interactions.${i}.options`, [
                          ...current,
                          { id: uuidv4(), label: "", role_native_id: "" },
                        ]);
                      }}
                    >
                      Thêm option
                    </Button>
                  </>
                )}

                <Button
                  size="xs"
                  color="red"
                  leftSection={<IconTrashFilled size={14} />}
                  onClick={() => removeInteraction(i)}
                >
                  Xóa interaction
                </Button>
              </Stack>
            </Fieldset>
          ))}
        </Stack>

        <Space h="lg" />

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
                <Flex className={classes.headerFooterFlexWrapper}>
                  <FileButton
                    onChange={(f) => {
                      if (f) {
                        console.log(f);
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

                  <TextInput
                    placeholder="Header"
                    className={classes.headerFooterTextInput}
                    {...form.getInputProps("headerMessage")}
                  />
                </Flex>

                <Group className={classes.titleAndMainContentArea}>
                  <Stack className={classes.titleAndMainText}>
                    <TextInput
                      placeholder="Title"
                      {...form.getInputProps("titleMessage")}
                    />
                    <Textarea
                      placeholder="Main content"
                      maxRows={9}
                      minRows={3}
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
                        content={{
                          allowedContent: "image/png,image/jpeg,image/webp",
                        }}
                        onBeforeUploadBegin={async (files) => {
                          return Promise.all(files.map(optimizeImage));
                        }}
                        onClientUploadComplete={(res) => {
                          setMainImageUrl(res[0].ufsUrl);
                          form.setFieldValue("featuredImage", res[0].ufsUrl);
                        }}
                        onUploadError={(error: Error) => {
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
                    <TextInput
                      placeholder="Field name"
                      {...form.getInputProps(`customFields.${i}.fieldName`)}
                    />
                    <TextInput
                      placeholder="Field value"
                      mt="md"
                      {...form.getInputProps(`customFields.${i}.fieldValue`)}
                    />
                    <Button
                      leftSection={<IconTrashFilled size={14} />}
                      color={vars.colors.red[9]}
                      mt="md"
                      onClick={() => {
                        form.removeListItem("customFields", i);
                      }}
                    >
                      Delete field
                    </Button>
                  </Fieldset>
                ))}
                {!(cFields.length >= MAX_CUSTOM_FIELDS) && (
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
                )}

                <Divider variant="dotted" my="md" />

                <Flex className={classes.headerFooterFlexWrapper}>
                  <FileButton
                    onChange={(f) => {
                      if (f) {
                        console.log(f);
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

                  <TextInput
                    placeholder="Footer"
                    className={classes.headerFooterTextInput}
                    {...form.getInputProps("footerMessage")}
                  />
                </Flex>
              </Stack>
            </Group>

            <Button
              variant="gradient"
              type="submit"
              loading={isPublishing}
              disabled={interactions.length === 0 || !form.values.channel_id}
            >
              Publish to Discord
            </Button>
          </Stack>
        </Group>
      </Stack>
    </form>
  );
};

export default EmbedEditor;
