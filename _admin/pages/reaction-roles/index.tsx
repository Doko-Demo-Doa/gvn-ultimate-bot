import React from "react";
import { GetServerSidePropsContext, InferGetStaticPropsType } from "next";

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
    <div>
      <div>{3}</div>
      <div>{4}</div>
    </div>
  );
};

ReactionRolesPage.getLayout = function getLayout(page: React.ReactNode) {
  return <>{page}</>;
};

export default ReactionRolesPage;
