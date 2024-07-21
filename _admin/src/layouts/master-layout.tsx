"use client";

import { AppShell, Group, NavLink, Title } from "@mantine/core";
import {
  IconGitPullRequest,
  IconBrandFramerMotion,
  IconHistory,
  IconRosetteFilled,
} from "@tabler/icons-react";
import Link from "next/link";
import { useDisclosure } from "@mantine/hooks";

import * as classes from "./master-layout.css";
import { usePathname } from "next/navigation";

interface Props {
  title?: string;
  description?: string;
  children?: React.ReactNode;
}

const MasterLayout: React.FC<Props> = ({ children, title, description }) => {
  const [opened, { toggle }] = useDisclosure();

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
}

function MainLink({ icon, color, label, to }: MainLinkProps) {
  const pathname = usePathname();

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
    icon: <IconBrandFramerMotion size="1rem" />,
    color: "teal",
    label: "Role Reaction Composer",
    to: "/reaction-roles",
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
  const links = data.map((link) => <MainLink {...link} key={link.label} />);
  return <div>{links}</div>;
}

export default MasterLayout;
