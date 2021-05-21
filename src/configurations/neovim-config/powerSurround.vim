nnoremap Q :call PowerSurround()<cr>

func! s:hasChar(inputChar)
  return match(getline("."), a:inputChar) > 0
endfunc

func! s:hasBoth(left, right)
  return s:hasChar(a:left) && s:hasChar(a:right)
endfunc

func! PowerSurround()
  if s:hasBoth('(', ')')
    exec "normal! vi("
  elseif s:hasBoth('[', ']')
    exec "normal! vi["
  elseif s:hasBoth('<', '>')
    exec "normal! vi<"
  elseif s:hasBoth('{', '}')
    exec "normal! vi{"
  endif
endfunc

