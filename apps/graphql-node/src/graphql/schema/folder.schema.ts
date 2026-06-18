export const folderTypeDefs = `#graphql
  type Folder {
    id: ID!
    name: String!
    parent_id: ID
    created_at: String!
  }

  extend type Query {
    folderTree: [Folder!]!
    folder(id: ID!): Folder
  }

  extend type Mutation {
    createFolder(name: String!, parentId: ID): Folder!
    updateFolder(id: ID!, name: String!): Folder!
    moveFolder(id: ID!, parentId: ID): Folder!
    deleteFolder(id: ID!): Boolean!
  }
`;
