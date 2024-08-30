import "@mantine/core/styles.css";
import "@mantine/notifications/styles.css";
import "@mantine/dropzone/styles.css";
import React from "react";
import { MantineProvider, ColorSchemeScript } from "@mantine/core";
import { Notifications } from "@mantine/notifications";
import { theme } from "../theme";
import { Metadata } from "next";
import { ClerkProvider } from "@clerk/nextjs";
import Providers from "~/providers/master-provider";
import { customApiClient } from "~/datasource/rest/api-client";

export const metadata: Metadata = {
  title: "DautoVN Admin",
  description: "Dashboard for DautoVN",
};

customApiClient.init({
  baseUrl: process.env.NEXT_PUBLIC_BASE_API_URL || "",
});

export default function RootLayout({ children }: { children: any }) {
  return (
    <ClerkProvider>
      <html lang="en">
        <head>
          <ColorSchemeScript />
          <link rel="shortcut icon" href="/favicon.svg" />
          <meta
            name="viewport"
            content="minimum-scale=1, initial-scale=1, width=device-width, user-scalable=no"
          />
        </head>
        <body>
          <Providers>
            <MantineProvider defaultColorScheme="dark" theme={theme}>
              {children}
              <Notifications />
            </MantineProvider>
          </Providers>
        </body>
      </html>
    </ClerkProvider>
  );
}
