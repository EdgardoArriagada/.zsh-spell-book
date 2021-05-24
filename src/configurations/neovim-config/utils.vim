func! Strip(input)
    return substitute(a:input, '^\s*\(.\{-}\)\s*$', '\1', '')
endfunc

func! IsEmptyLine(line)
    return len(Strip(getline(a:line))) == 0
endfunc
