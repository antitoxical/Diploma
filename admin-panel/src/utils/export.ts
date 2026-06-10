import * as XLSX from 'xlsx';
import { saveAs } from 'file-saver';

export const exportToExcel = (data: any[], filename: string) => {
  const worksheet = XLSX.utils.json_to_sheet(data);
  const workbook = XLSX.utils.book_new();
  XLSX.utils.book_append_sheet(workbook, worksheet, 'Data');
  
  const excelBuffer = XLSX.write(workbook, { bookType: 'xlsx', type: 'array' });
  const blob = new Blob([excelBuffer], { type: 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet' });
  saveAs(blob, `${filename}.xlsx`);
};

export const exportToCSV = (data: any[], filename: string) => {
  const worksheet = XLSX.utils.json_to_sheet(data);
  const csv = XLSX.utils.sheet_to_csv(worksheet);
  
  const blob = new Blob([csv], { type: 'text/csv;charset=utf-8;' });
  saveAs(blob, `${filename}.csv`);
};

// Helper to flatten nested objects for export
export const flattenObject = (obj: any, prefix = ''): any => {
  return Object.keys(obj).reduce((acc: any, key: string) => {
    const pre = prefix.length ? `${prefix}.` : '';
    if (typeof obj[key] === 'object' && obj[key] !== null && !Array.isArray(obj[key])) {
      Object.assign(acc, flattenObject(obj[key], pre + key));
    } else if (Array.isArray(obj[key])) {
      acc[pre + key] = obj[key].join(', ');
    } else {
      acc[pre + key] = obj[key];
    }
    return acc;
  }, {});
};
