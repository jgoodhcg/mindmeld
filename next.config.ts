import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  reactCompiler: true,
  turbopack: {
    root: process.cwd(), // fixes root inference when multiple lockfiles exist
  },
};

export default nextConfig;
