import { createVanillaExtractPlugin } from "@vanilla-extract/next-plugin";


const withVanillaExtract = createVanillaExtractPlugin({
  unstable_turbopack: { mode: "auto" },
});

/** @type {import('next').NextConfig} */
const nextConfig = {
  reactStrictMode: true,
  experimental: {
    optimizePackageImports: ["@mantine/core", "@mantine/hooks"],
  },
  turbopack: {},
};

export default withVanillaExtract(nextConfig);
