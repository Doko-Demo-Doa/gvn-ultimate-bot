import {
  AppShell,
  Header,
  Text,
  Navbar,
  ThemeIcon,
  UnstyledButton,
  Group,
  Title,
} from "@mantine/core";
import {
  IconGitPullRequest,
  IconBrandFramerMotion,
  IconHistory,
  IconRosetteFilled,
} from "@tabler/icons-react";
import NextLink from "next/link";

interface Props {
  title?: string;
  description?: string;
  children?: React.ReactNode;
}

const MasterLayout: React.FC<Props> = ({ children, title, description }) => {
  return (
    <AppShell
      padding={0}
      header={
        <Header height={60} px={20} p="xs">
          <Title order={2}>DautoVN Control Panel</Title>
        </Header>
      }
      navbar={
        <Navbar p="md" hiddenBreakpoint="sm" width={{ sm: 200, lg: 300 }}>
          <Navbar.Section grow mt="md">
            <MainLinks />
          </Navbar.Section>
        </Navbar>
      }
    >
      {children}
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
  return (
    <UnstyledButton
      component={NextLink}
      href={to}
      sx={(theme) => ({
        display: "block",
        width: "100%",
        padding: theme.spacing.xs,
        borderRadius: theme.radius.sm,
        color:
          theme.colorScheme === "dark" ? theme.colors.dark[0] : theme.black,

        "&:hover": {
          backgroundColor:
            theme.colorScheme === "dark"
              ? theme.colors.dark[6]
              : theme.colors.gray[0],
        },
      })}
    >
      <Group>
        <ThemeIcon color={color} variant="light">
          {icon}
        </ThemeIcon>

        <Text size="sm">{label}</Text>
      </Group>
    </UnstyledButton>
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
    to: "/role-reaction",
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
