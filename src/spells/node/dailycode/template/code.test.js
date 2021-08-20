function getResult(input) {
  return Object.entries(input).map(([k, v]) => ({[k]: v}))
}

describe('daily code problem', () => {
  it('should succesfully test', () => {
    const input = {dos: 2, uno: 1, tres: 3}

    const result = getResult(input)

    expect(result).toEqual([{dos: 2}, {uno: 1}, {tres: 3}])
  })
})
