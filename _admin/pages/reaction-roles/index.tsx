import { GetServerSidePropsContext, InferGetStaticPropsType } from "next";
import { Box, Container, Divider, Paper, Title } from "@mantine/core";
import MasterLayout from "~/layouts/master-layout";
import EmbedEditor from "~/components/embed-editor/embed-editor";

export async function getServerSideProps(context: GetServerSidePropsContext) {
  return {
    props: {},
  };
}

type Props = React.FC<InferGetStaticPropsType<typeof getServerSideProps>> & {
  getLayout: (page: React.ReactNode) => JSX.Element;
};

const ReactionRolesPage: Props = ({}) => {
  return (
    <Box p="lg">
      <Paper>
        <Title order={3}>Your reaction role messages</Title>
        <Divider my="sm" />

        <EmbedEditor messageId="" />
      </Paper>
    </Box>
  );
};

ReactionRolesPage.getLayout = function getLayout(page: React.ReactNode) {
  return (
    <MasterLayout title="Reaction Roles" description="CP">
      {page}
    </MasterLayout>
  );
};

export default ReactionRolesPage;
