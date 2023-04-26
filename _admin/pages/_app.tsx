import Head from "next/head";
import { AppContext, AppInitialProps, AppLayoutProps } from "next/app";
import { NextComponentType } from "next";

import { MantineProvider } from "@mantine/core";
import { SessionProvider } from "next-auth/react";
import type { Session } from "next-auth";
import {
  useQuery,
  useMutation,
  useQueryClient,
  QueryClient,
  QueryClientProvider,
} from "@tanstack/react-query";
const queryClient = new QueryClient();

const App: NextComponentType<AppContext, AppInitialProps, AppLayoutProps> = ({
  Component,
  pageProps: { session, ...pageProps },
}: AppLayoutProps<{ session: Session }>) => {
  const getLayout = Component.getLayout || ((page: React.ReactNode) => page);

  return (
    <>
      <Head>
        <title>DautoVN</title>
        <meta
          name="viewport"
          content="minimum-scale=1, initial-scale=1, width=device-width"
        />
      </Head>
      <QueryClientProvider client={queryClient}>
        <SessionProvider session={session}>
          <MantineProvider
            withGlobalStyles
            withNormalizeCSS
            theme={{
              colorScheme: "dark",
            }}
          >
            {getLayout(<Component {...pageProps} />)}
          </MantineProvider>
        </SessionProvider>
      </QueryClientProvider>
    </>
  );
};

export default App;
