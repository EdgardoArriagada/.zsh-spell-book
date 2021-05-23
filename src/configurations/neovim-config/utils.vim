func! Strip(input)
    return substitute(a:input, '^\s*\(.\{-}\)\s*$', '\1', '')
endfunc

func! IsEmptyLine()
    return Strip(len(getline("."))) == 0
endfunc

