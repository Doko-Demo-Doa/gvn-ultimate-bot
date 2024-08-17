"use client";

import { Button, Group, NumberInput, Text } from "@mantine/core";
import MasterLayout from "~/layouts/master-layout";

export default function Page() {
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
