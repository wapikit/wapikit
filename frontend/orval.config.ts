import { defineConfig } from 'orval';

export default defineConfig({
  api: {
    input: '../swagger/collections.yaml',
    output: {
      packageJson: './package.json',
      mode: 'single',
      prettier: true,
      client: 'react-query',
      tsconfig: './tsconfig.json',
      target: './.generated.ts',
    },
    hooks: {
      afterAllFilesWrite: "prettier --write"
    }
  },
});