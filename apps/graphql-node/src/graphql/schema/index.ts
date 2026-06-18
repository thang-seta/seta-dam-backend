import { folderTypeDefs } from './folder.schema';
import { metadataTypeDefs } from './metadata.schema';
import { permissionTypeDefs } from './permission.schema';

const baseTypeDefs = `#graphql
  type Query {
    _empty: String
  }
  type Mutation {
    _empty: String
  }
`;

export const typeDefs = [
  baseTypeDefs,
  folderTypeDefs,
  metadataTypeDefs,
  permissionTypeDefs,
];
