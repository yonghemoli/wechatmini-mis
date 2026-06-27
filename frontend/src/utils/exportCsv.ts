export const exportCsv = (
  filename: string,
  rows: Array<Record<string, string | number | boolean>>
) => {
  if (!rows.length) return
  const headers = Object.keys(rows[0])
  const escape = (value: string | number | boolean) =>
    `"${String(value).replace(/"/g, '""')}"`
  const csv = [
    headers.join(','),
    ...rows.map(row => headers.map(header => escape(row[header])).join(','))
  ].join('\n')
  const blob = new Blob([`\uFEFF${csv}`], {
    type: 'text/csv;charset=utf-8;'
  })
  const url = URL.createObjectURL(blob)
  const link = document.createElement('a')
  link.href = url
  link.download = filename
  link.click()
  URL.revokeObjectURL(url)
}
