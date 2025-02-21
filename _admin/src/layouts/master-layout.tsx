"use client";

import { AppShell, Group, Loader, NavLink, Title } from "@mantine/core";
import { useDisclosure } from "@mantine/hooks";
import {
  IconBrandFramerMotion,
  IconGitPullRequest,
  IconHistory,
  IconPin,
  IconRosetteFilled,
} from "@tabler/icons-react";
import Link from "next/link";

import { usePathname } from "next/navigation";
import { BotModuleConst } from "~/common/bot-module-const";
import { customApiClient } from "~/datasource/rest/api-client";
import { useAppModules } from "~/hooks/api-hooks";
import * as classes from "./master-layout.css";

import "@uploadthing/react/styles.css";

interface Props {
  title?: string;
  description?: string;
  children?: React.ReactNode;
}

// We have to init twice, one for client and one for server
customApiClient.init({
  baseUrl: process.env.NEXT_PUBLIC_BASE_API_URL || "",
});

const MasterLayout: React.FC<Props> = ({ children }) => {
  const [opened] = useDisclosure();

  return (
    <AppShell
      header={{ height: 60 }}
      navbar={{ width: 300, breakpoint: "sm", collapsed: { mobile: !opened } }}
      padding="md"
    >
      <AppShell.Header>
        <Group className={classes.header}>
          <Title order={3}>DautoVN Bot CP</Title>
        </Group>
      </AppShell.Header>

      <AppShell.Navbar px="lg" py="lg">
        <MainLinks />
      </AppShell.Navbar>

      <AppShell.Main>{children}</AppShell.Main>
    </AppShell>
  );
};

interface MainLinkProps {
  icon: React.ReactNode;
  color: string;
  label: string;
  to: string;
  disabled?: boolean;
  internalModuleName?: string;
}

function MainLink({ icon, color, label, disabled, to }: MainLinkProps) {
  const pathname = usePathname();

  if (disabled) {
    return null;
  }

  return (
    <NavLink
      component={Link}
      href={to}
      className={classes.menuItem}
      label={label}
      active={to === pathname}
      color={color}
      leftSection={icon}
    />
  );
}

const data = [
  {
    icon: <IconGitPullRequest size="1rem" />,
    color: "blue",
    label: "Module Enabler",
    to: "/",
  },
  {
    icon: <IconPin size="1rem" />,
    color: "orange",
    label: "Pin module",
    to: "/pin",
    internalModuleName: BotModuleConst.PIN_MODULE,
  },
  {
    icon: <IconBrandFramerMotion size="1rem" />,
    color: "teal",
    label: "Role Reaction Composer",
    to: "/reaction-roles",
    internalModuleName: BotModuleConst.REACTION_ROLE_MODULE,
  },
  {
    icon: <IconHistory size="1rem" />,
    color: "violet",
    label: "Server Message Log",
    to: "/server-log",
  },
  {
    icon: <IconRosetteFilled size="1rem" />,
    color: "grape",
    label: "Role Manager",
    to: "/role-manager",
  },
];

export function MainLinks() {
  const { data: remoteData } = useAppModules();

  if (!remoteData) {
    return <Loader color="teal" />;
  }

  const links = data.map((link) => (
    <MainLink
      {...link}
      disabled={
        remoteData.data.find((n) => n.ModuleName === link.internalModuleName)
          ?.IsActivated === 0
      }
      key={link.label}
    />
  ));
  return <div>{links}</div>;
}

export default MasterLayout;
