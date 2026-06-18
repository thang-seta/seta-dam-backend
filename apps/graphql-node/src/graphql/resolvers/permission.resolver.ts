import { GoServiceClient } from '../../clients/go-service.client';

function validateObjectType(objectType: string): string {
  if (objectType === 'folder' || objectType === 'metadata' || objectType === 'metadata_item') {
    return objectType;
  }
  throw new Error('objectType must be folder, metadata, or metadata_item');
}

function validateAction(action: string): string {
  if (action === 'read' || action === 'write' || action === 'manage_permissions') {
    return action;
  }
  throw new Error('action must be read, write, or manage_permissions');
}

export const permissionResolvers = {
  Query: {
    effectivePermissions: async (
      _: any,
      { userId, objectType, objectId }: { userId: string; objectType: string; objectId: string },
      context: any
    ) => {
      return GoServiceClient.getEffectivePermissions(userId, validateObjectType(objectType), objectId, context.user);
    },
  },
  Mutation: {
    grantPermission: async (
      _: any,
      { userId, objectType, objectId, action }: { userId: string; objectType: string; objectId: string; action: string },
      context: any
    ) => {
      const payload = {
        user_id: userId,
        object_type: validateObjectType(objectType),
        object_id: objectId,
        action: validateAction(action),
      };
      await GoServiceClient.grantPermission(payload, context.user);
      return true;
    },
    revokePermission: async (
      _: any,
      { userId, objectType, objectId, action }: { userId: string; objectType: string; objectId: string; action: string },
      context: any
    ) => {
      const payload = {
        user_id: userId,
        object_type: validateObjectType(objectType),
        object_id: objectId,
        action: validateAction(action),
      };
      await GoServiceClient.revokePermission(payload, context.user);
      return true;
    },
  },
};
