export const permissionTypeDefs = `#graphql
  type EffectivePermissions {
    user_id: ID!
    object_type: String!
    object_id: ID!
    actions: [String!]!
  }

  extend type Query {
    effectivePermissions(userId: ID!, objectType: String!, objectId: ID!): EffectivePermissions!
  }

  extend type Mutation {
    grantPermission(userId: ID!, objectType: String!, objectId: ID!, action: String!): Boolean!
    revokePermission(userId: ID!, objectType: String!, objectId: ID!, action: String!): Boolean!
  }
`;
