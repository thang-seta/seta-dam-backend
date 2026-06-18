import { GoServiceClient } from '../../clients/go-service.client';

export const folderResolvers = {
  Query: {
    folderTree: async (_: any, __: any, context: any) => {
      return GoServiceClient.getFolderTree(context.user);
    },
    folder: async (_: any, { id }: { id: string }, context: any) => {
      return GoServiceClient.getFolderById(id, context.user);
    },
  },
  Mutation: {
    createFolder: async (_: any, { name, parentId }: { name: string; parentId?: string | null }, context: any) => {
      return GoServiceClient.createFolder(name, parentId || null, context.user);
    },
    updateFolder: async (_: any, { id, name }: { id: string; name: string }, context: any) => {
      await GoServiceClient.updateFolder(id, name, context.user);
      return GoServiceClient.getFolderById(id, context.user);
    },
    moveFolder: async (_: any, { id, parentId }: { id: string; parentId?: string | null }, context: any) => {
      await GoServiceClient.moveFolder(id, parentId || null, context.user);
      return GoServiceClient.getFolderById(id, context.user);
    },
    deleteFolder: async (_: any, { id }: { id: string }, context: any) => {
      await GoServiceClient.deleteFolder(id, context.user);
      return true;
    },
  },
};
