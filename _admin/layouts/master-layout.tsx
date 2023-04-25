import { AppShell, Header, Text, Navbar } from "@mantine/core";

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
        <Header height={60} p="xs">
          DautoVNs
        </Header>
      }
      navbar={
        <Navbar p="md" hiddenBreakpoint="sm" width={{ sm: 200, lg: 300 }}>
          <Text>Application navbar</Text>
        </Navbar>
      }
    >
      {children}
    </AppShell>
  );
};

export default MasterLayout;
