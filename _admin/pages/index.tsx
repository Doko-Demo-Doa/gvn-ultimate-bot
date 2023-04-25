import React from "react";
import { Text } from "@mantine/core";
import MasterLayout from "~/layouts/master-layout";

type Props = React.FC<{}> & {
  getLayout: (page: React.ReactNode) => JSX.Element;
};

const HomepageRoute: Props = () => {
  return (
    <>
      <Text>xxxxx</Text>
    </>
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
