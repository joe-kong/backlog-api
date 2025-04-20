import { ApolloClient, InMemoryCache, createHttpLink } from '@apollo/client';
import { setContext } from '@apollo/client/link/context';

// REST APIをGraphQLで扱うための設定
const httpLink = createHttpLink({
  uri: '/graphql', // バックエンドのGraphQLエンドポイント（実際にはREST APIをマッピング）
});

// 認証用のヘッダーを追加
const authLink = setContext((_, { headers }) => {
  const token = localStorage.getItem('accessToken');
  return {
    headers: {
      ...headers,
      authorization: token ? `Bearer ${token}` : '',
    },
  };
});

// Apollo Clientの初期化
export const client = new ApolloClient({
  link: authLink.concat(httpLink),
  cache: new InMemoryCache(),
}); 