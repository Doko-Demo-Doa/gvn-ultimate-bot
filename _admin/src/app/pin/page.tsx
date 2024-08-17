"use client";

import { Button, Group, NumberInput, Text } from "@mantine/core";
import { useAppModule } from "~/hooks/api-hooks";
import MasterLayout from "~/layouts/master-layout";

export default function Page() {
  const PIN_ID = 1;
  const { data } = useAppModule(PIN_ID);
  console.log("dd", data?.data);

  return (
    <MasterLayout>
      <Group>
        <NumberInput placeholder="Max threshold" />
        <Button>Save</Button>
        <Text>Saved</Text>
      </Group>
    </MasterLayout>
  );
}
