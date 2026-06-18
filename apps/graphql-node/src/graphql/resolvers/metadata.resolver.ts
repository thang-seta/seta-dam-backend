import { GoServiceClient } from '../../clients/go-service.client';

export const metadataResolvers = {
  Query: {
    metadataList: async (_: any, { folderId }: { folderId: string }, context: any) => {
      const res = await GoServiceClient.getMetadataList(folderId, context.user);
      return res || [];
    },
    metadataDetail: async (_: any, { id }: { id: string }, context: any) => {
      return GoServiceClient.getMetadataById(id, context.user);
    },
  },
  Mutation: {
    createMetadata: async (
      _: any,
      args: {
        folderId: string;
        title: string;
        description?: string;
        labels?: string[];
        category?: string;
        externalSource?: string;
        externalId?: string;
        sourceUrl?: string;
        thumbnailUrl?: string;
        license?: string;
        author?: string;
        metadataJson?: string;
        notes?: string;
      },
      context: any
    ) => {
      const payload = {
        folder_id: args.folderId,
        title: args.title,
        description: args.description || '',
        labels: args.labels || [],
        category: args.category || '',
        external_source: args.externalSource || '',
        external_id: args.externalId || '',
        source_url: args.sourceUrl || '',
        thumbnail_url: args.thumbnailUrl || '',
        license: args.license || '',
        author: args.author || '',
        metadata_json: args.metadataJson || '',
        notes: args.notes || '',
      };
      return GoServiceClient.createMetadata(payload, context.user);
    },
    updateMetadata: async (
      _: any,
      args: {
        id: string;
        folderId: string;
        title: string;
        description?: string;
        labels?: string[];
        category?: string;
        externalSource?: string;
        externalId?: string;
        sourceUrl?: string;
        thumbnailUrl?: string;
        license?: string;
        author?: string;
        metadataJson?: string;
        notes?: string;
      },
      context: any
    ) => {
      const payload = {
        id: args.id,
        folder_id: args.folderId,
        title: args.title,
        description: args.description || '',
        labels: args.labels || [],
        category: args.category || '',
        external_source: args.externalSource || '',
        external_id: args.externalId || '',
        source_url: args.sourceUrl || '',
        thumbnail_url: args.thumbnailUrl || '',
        license: args.license || '',
        author: args.author || '',
        metadata_json: args.metadataJson || '',
        notes: args.notes || '',
      };
      await GoServiceClient.updateMetadata(args.id, payload, context.user);
      return GoServiceClient.getMetadataById(args.id, context.user);
    },
    deleteMetadata: async (_: any, { id }: { id: string }, context: any) => {
      await GoServiceClient.deleteMetadata(id, context.user);
      return true;
    },
  },
};
