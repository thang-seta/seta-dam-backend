export const metadataTypeDefs = `#graphql
  type Metadata {
    id: ID!
    folder_id: ID!
    title: String!
    description: String
    labels: [String!]
    category: String
    source_url: String
    notes: String
    created_at: String!
  }

  extend type Query {
    metadataList(folderId: ID!): [Metadata!]!
    metadataDetail(id: ID!): Metadata
  }

  extend type Mutation {
    createMetadata(
      folderId: ID!
      title: String!
      description: String
      labels: [String!]
      category: String
      sourceUrl: String
      notes: String
    ): Metadata!

    updateMetadata(
      id: ID!
      folderId: ID!
      title: String!
      description: String
      labels: [String!]
      category: String
      sourceUrl: String
      notes: String
    ): Metadata!

    deleteMetadata(id: ID!): Boolean!
  }
`;
