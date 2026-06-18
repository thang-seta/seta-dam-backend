export const metadataTypeDefs = `#graphql
  type Metadata {
    id: ID!
    folder_id: ID!
    title: String!
    description: String
    labels: [String!]
    category: String
    external_source: String
    external_id: String
    source_url: String
    thumbnail_url: String
    license: String
    author: String
    metadata_json: String
    notes: String
    created_by: ID!
    updated_by: ID
    created_at: String!
    updated_at: String!
  }

  extend type Query {
    metadataItems: [Metadata!]!
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
      externalSource: String
      externalId: String
      sourceUrl: String
      thumbnailUrl: String
      license: String
      author: String
      metadataJson: String
      notes: String
    ): Metadata!

    updateMetadata(
      id: ID!
      folderId: ID!
      title: String!
      description: String
      labels: [String!]
      category: String
      externalSource: String
      externalId: String
      sourceUrl: String
      thumbnailUrl: String
      license: String
      author: String
      metadataJson: String
      notes: String
    ): Metadata!

    deleteMetadata(id: ID!): Boolean!
  }
`;
