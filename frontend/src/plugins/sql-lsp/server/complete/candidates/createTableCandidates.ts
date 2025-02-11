import { uniqBy } from "lodash-es";
import { CompletionItem } from "vscode-languageserver-types";
import { Table } from "../../../types";
import { ICONS, SortText } from "../utils";

export const createTableCandidates = (
  tableList: Table[],
  withDatabasePrefix = true
): CompletionItem[] => {
  const suggestions: CompletionItem[] = [];

  tableList.forEach((table) => {
    const label = withDatabasePrefix
      ? `${table.database}.${table.name}`
      : table.name;
    suggestions.push({
      label,
      kind: ICONS.TABLE,
      detail: "<Table>",
      sortText: SortText.TABLE,
      insertText: label,
    });
  });

  return uniqBy(suggestions, (item) => item.label);
};
