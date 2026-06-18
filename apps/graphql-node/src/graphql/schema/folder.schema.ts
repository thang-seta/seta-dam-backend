export const folderTypeDefs = `#graphql
  type Folder {
    id: ID!
    name: String!
    parent_id: ID
    description: String
    created_by: ID!
    updated_by: ID
    created_at: String!
    updated_at: String!
  }

  extend type Query {
    folderTree: [Folder!]!
    folder(id: ID!): Folder
  }

  extend type Mutation {
    createFolder(name: String!, description: String, parentId: ID): Folder!
    updateFolder(id: ID!, name: String!, description: String): Folder!
    moveFolder(id: ID!, parentId: ID): Folder!
    deleteFolder(id: ID!): Boolean!
  }
`;
