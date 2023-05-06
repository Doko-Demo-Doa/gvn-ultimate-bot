import React from "react";
import { Container, Loader, Stack, Switch, Title } from "@mantine/core";
import axios from "axios";

import MasterLayout from "~/layouts/master-layout";
import { useAppModuleEnabler, useAppModules } from "~/hooks/use-app-modules";

type Props = React.FC<{}> & {
  getLayout: (page: React.ReactNode) => JSX.Element;
};

const LABEL_MAPPER: Record<string, string> = {
  pin_module: "Tự động pin nếu đủ 4 reaction",
  grant_role_module: "React tạo role mạnh mẽ",
};

const HomepageRoute: Props = () => {
  const { data, isLoading } = useAppModules();

  console.log(data);
  const { mutate } = useAppModuleEnabler();

  if (!data || isLoading) {
    return <Loader />;
  }

  return (
    <Stack p="lg">
      <Title
        onClick={async () => {
          const BASE_URL = process.env.NEXT_PUBLIC_BACKEND_BASE_URL;

          const resp = await axios.get(BASE_URL + "/module/list");
          console.log(resp);
        }}
      >
        Module enabler
      </Title>

      {data.data.map((n) => (
        <Switch
          key={n.ID}
          label={LABEL_MAPPER[n.ModuleName]}
          value={n.IsActivated}
          onChange={(newVal) => {
            mutate(newVal.currentTarget.checked ? 1 : 0);
          }}
        />
      ))}
    </Stack>
  );
};

HomepageRoute.getLayout = function getLayout(page: React.ReactNode) {
  return (
    <MasterLayout title="DautoVN Control Panel" description="CP">
      {page}
    </MasterLayout>
  );
};

export default HomepageRoute;
