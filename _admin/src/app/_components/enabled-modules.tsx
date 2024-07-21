"use client";

import { Loader, Space, Stack, Switch, Title } from "@mantine/core";
import { useAppModuleEnabler, useAppModules } from "~/hooks/use-app-modules";

const moduleNameMap: Record<string, string> = {
  pin_module: "Pin Module",
  grant_role_module: "Grant Role Module",
};

const EnabledModules = () => {
  const { data } = useAppModules();
  const { mutate } = useAppModuleEnabler();

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
            onChange={(event) =>
              mutate({
                id: module.ID,
                activated: event.target.checked,
              })
            }
          />
        ))}
      </Stack>
    </>
  );
};

export default EnabledModules;
