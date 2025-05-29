- Start Date: 2025-04-25
- Target Version: 1.0

# Summary

Currently, our installed Material-UI library is at v4, but the latest version, v7, offers significant improvements and benefits.

# Motivation

Material-UI v4 is no longer actively maintained, making it harder to address bugs, security issues, or compatibility with modern tools.

More over, Upgrading to Material-UI v7 provides the following benefits:

- **Improved Performance**: Optimized components and reduced bundle size for faster rendering.
- **Modern Features**: Access to new components and hooks introduced in v5, v6, and v7.
- **Long-Term Support**: Continued updates, bug fixes, and security patches for the latest version.
- **Improved Theming**: More flexible and powerful theming capabilities.
- **Compatibility**: Align with the latest React versions and ecosystem tools.
- **Community and Ecosystem**: The community and ecosystem have largely moved to newer versions, making it easier to find resources.

# Detailed design

To upgrade from Material-UI v4 to v7, the following steps will be taken:

## Step 1: Upgrade from v4 to v5

1. **Upgrade Dependencies to v5**:

   - Make sure React version 17.0.0 and above
   - Update typescript

2. **Resolve Breaking Changes**:

   - Review the [v4 to v5 migration guide](https://mui.com/material-ui/migration/migration-v4/) for breaking changes and migration instruction
   - Update imports, component props, and theming configurations as per the migration guide.
   - Run codemods using: 

   ```bash
   npx @mui/codemod@latest v5.0.0/preset-safe <path>
   ```

3. **Update Theming**:

   - Replace the legacy `createMuiTheme` with `createTheme`.
   - Update theme configuration to align with the new API.

4. **Test and Validate**

## Step 2: Upgrade from v5 to v6

1. **Upgrade Dependencies to v6**:

2. **Resolve Breaking Changes**:

   - Review the [v5 to v6 migration guide](https://mui.com/material-ui/migration/migration-v5/) for breaking changes.
   - Update imports, component props, and theming configurations as per the migration guide.

3. **Test and Validate**

## Step 3: Upgrade from v6 to v7

1. **Upgrade Dependencies to v7**:

2. **Resolve Breaking Changes**:

   - Review the [v6 to v7 migration guide](https://mui.com/material-ui/migration/migration-v6/) for breaking changes.
   - Update imports, component props, and theming configurations as per the migration guide.

3. **Test and Validate**