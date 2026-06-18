import { GoServiceClient } from '../../clients/go-service.client';

export const permissionResolvers = {
  Query: {
    effectivePermissions: async (
      _: any,
      { userId, objectType, objectId }: { userId: string; objectType: string; objectId: string },
      context: any
    ) => {
      return GoServiceClient.getEffectivePermissions(userId, objectType, objectId, context.user);
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
        object_type: objectType,
        object_id: objectId,
        action,
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
        object_type: objectType,
        object_id: objectId,
        action,
      };
      await GoServiceClient.revokePermission(payload, context.user);
      return true;
    },
  },
};
