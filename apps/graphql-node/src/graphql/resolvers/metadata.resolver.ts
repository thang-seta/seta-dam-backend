import { GoServiceClient } from '../../clients/go-service.client';

function requireNonBlank(value: string, fieldName: string): string {
  const normalized = value.trim();
  if (!normalized) {
    throw new Error(`${fieldName} is required`);
  }
  return normalized;
}

function validateMetadataJson(value?: string): string {
  const normalized = value?.trim() || '';
  if (!normalized) {
    return '';
  }
  try {
    JSON.parse(normalized);
    return normalized;
  } catch {
    throw new Error('metadataJson must be valid JSON');
  }
}

export const metadataResolvers = {
  Query: {
    metadataItems: async (_: any, __: any, context: any) => {
      const res = await GoServiceClient.getAllMetadata(context.user);
      return res || [];
    },
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
        folder_id: requireNonBlank(args.folderId, 'folderId'),
        title: requireNonBlank(args.title, 'title'),
        description: args.description || '',
        labels: args.labels || [],
        category: args.category || '',
        external_source: args.externalSource || '',
        external_id: args.externalId || '',
        source_url: args.sourceUrl || '',
        thumbnail_url: args.thumbnailUrl || '',
        license: args.license || '',
        author: args.author || '',
        metadata_json: validateMetadataJson(args.metadataJson),
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
        folder_id: requireNonBlank(args.folderId, 'folderId'),
        title: requireNonBlank(args.title, 'title'),
        description: args.description || '',
        labels: args.labels || [],
        category: args.category || '',
        external_source: args.externalSource || '',
        external_id: args.externalId || '',
        source_url: args.sourceUrl || '',
        thumbnail_url: args.thumbnailUrl || '',
        license: args.license || '',
        author: args.author || '',
        metadata_json: validateMetadataJson(args.metadataJson),
        notes: args.notes || '',
      };
      return GoServiceClient.updateMetadata(args.id, payload, context.user);
    },
    deleteMetadata: async (_: any, { id }: { id: string }, context: any) => {
      await GoServiceClient.deleteMetadata(id, context.user);
      return true;
    },
  },
};
