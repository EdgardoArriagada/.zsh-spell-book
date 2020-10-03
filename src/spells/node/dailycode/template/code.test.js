function solverFunction (array) {
    return array
}
describe('daily code problem', () => {
    it('should succesfully test', () => {
	const input = [0, 1, 2]
	const expectedResult = [0, 1, 2]    
        expect(solverFunction(input)).toEqual(expectedResult)
    })
})
