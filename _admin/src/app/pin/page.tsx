"use client";

import { Button, Group, Loader, NumberInput, Stack, Text } from "@mantine/core";
import { useRef } from "react";
import { notifications } from "@mantine/notifications";
import { useAppModule, useModuleConfigMutation } from "~/hooks/api-hooks";
import MasterLayout from "~/layouts/master-layout";
import { IPinModuleConfig } from "~/types/payload";

const PIN_ID = 1;

function tryParse(input: string) {
  try {
    return JSON.parse(input) as IPinModuleConfig;
  } catch (e) {
    return null;
  }
}

export default function Page() {
  const { data, isLoading } = useAppModule(PIN_ID);
  const { mutateAsync, isPending } = useModuleConfigMutation();

  let defaultValue = 0;
  if (data) {
    const parsed = tryParse(data.data.CustomConfig);

    if (parsed !== null) {
      defaultValue = parsed.threshold;
    }
  }

  const val = useRef(defaultValue);

  return (
    <MasterLayout>
      {isLoading ? (
        <Loader />
      ) : (
        <Stack>
          <Text>
            Số lượng dấu Pin reaction tối thiểu để có thể pin message lên. Nếu
            dưới số này thì message sẽ bị unpin.
          </Text>

          <Group>
            <NumberInput
              defaultValue={defaultValue}
              onChange={(v) => {
                val.current = Number(v);
              }}
              placeholder="Max threshold"
            />
            <Button
              loading={isPending}
              onClick={async () => {
                const newConfig: { threshold: number } = {
                  threshold: val.current,
                };

                const resp = await mutateAsync({
                  module_id: PIN_ID,
                  new_config: newConfig,
                });

                if (resp) {
                  notifications.show({
                    title: "Saved",
                    message: `Con số ngưỡng mới là ${val.current}`,
                  });
                }
              }}
            >
              Save
            </Button>
          </Group>
        </Stack>
      )}
    </MasterLayout>
  );
}
