import React from "ract";
import { NextPage, NextPageContext } from "next";
import Error from "next/error";

interface InitialProps {
  exampleProp?: any;
  statusCode?: number;
}

const TemplateComponent: NextPage<InitialProps> = ({
  exampleProp,
  statusCode,
}) => {
  if (statusCode) return <Error statusCode={statusCode} />;

  return (
    <>
      <div>{`TemplateComponent Component ${exampleProp}`}</div>
    </>
  );
};

TemplateComponent.getInitialProps = async (
  context: NextPageContext
): Promise<InitialProps> => {
  const { query } = context;
  try {
    // const exampleQueryParam: string = normalizeQueryParam(query.exampleQueryParam);

    //const url = `https://google.cl/search?q={exampleQueryParam}`;

    //const response = await Axios.get(url);
    return {
      //  exampleProp: response?.data,
      exampleProp: "EXAMPLE PROP",
    };
  } catch (error) {
    return {
      statusCode: error?.status ?? error?.response?.status ?? 500,
    };
  }
};

export default TemplateComponent;
