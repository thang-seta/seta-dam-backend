import { GoServiceClient } from '../../clients/go-service.client';

export const folderResolvers = {
  Query: {
    folderTree: async (_: any, __: any, context: any) => {
      const res = await GoServiceClient.getFolderTree(context.user);
      return res || [];
    },
    folder: async (_: any, { id }: { id: string }, context: any) => {
      return GoServiceClient.getFolderById(id, context.user);
    },
  },
  Mutation: {
    createFolder: async (
      _: any,
      { name, description, parentId }: { name: string; description?: string | null; parentId?: string | null },
      context: any
    ) => {
      return GoServiceClient.createFolder(name, description || null, parentId || null, context.user);
    },
    updateFolder: async (
      _: any,
      { id, name, description }: { id: string; name: string; description?: string | null },
      context: any
    ) => {
      await GoServiceClient.updateFolder(id, name, description || null, context.user);
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
