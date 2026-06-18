import { folderResolvers } from './folder.resolver';
import { metadataResolvers } from './metadata.resolver';
import { permissionResolvers } from './permission.resolver';

export const resolvers = {
  Query: {
    ...folderResolvers.Query,
    ...metadataResolvers.Query,
    ...permissionResolvers.Query,
  },
  Mutation: {
    ...folderResolvers.Mutation,
    ...metadataResolvers.Mutation,
    ...permissionResolvers.Mutation,
  },
};
