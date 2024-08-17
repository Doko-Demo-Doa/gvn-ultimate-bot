"use client";

import { Loader, Space, Stack, Switch, Title } from "@mantine/core";
import {
  ModuleActivationStatus,
  useAppModuleEnabler,
  useAppModules,
} from "~/hooks/api-hooks";

const EnabledModules = () => {
  const { data, refetch } = useAppModules();
  const { mutateAsync } = useAppModuleEnabler();

  if (!data) {
    return <Loader />;
  }

  return (
    <>
      <Title order={3}>Enabled Modules</Title>
      <Space h="lg" />
      <Stack>
        {data?.data.map((module) => (
          <Switch
            key={module.ID}
            defaultChecked={!!module.IsActivated}
            label={module.ModuleLabel || ""}
            onChange={async (event) => {
              await mutateAsync({
                module_id: module.ID,
                is_activated: event.target.checked
                  ? ModuleActivationStatus.ENABLED
                  : ModuleActivationStatus.DISABLED,
              });

              refetch();
            }}
          />
        ))}
      </Stack>
    </>
  );
};

export default EnabledModules;
