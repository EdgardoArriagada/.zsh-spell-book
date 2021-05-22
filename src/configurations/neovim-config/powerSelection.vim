nnoremap Q :call PowerSelection()<cr>

func! s:hasChar(inputChar)
  return match(getline("."), a:inputChar) > 0
endfunc

func! s:hasPair(pair)
  return s:hasChar(a:pair[0]) && s:hasChar(a:pair[1])
endfunc

func! s:didSelectInline()
  return col('.') != col('v') && line('.') == line('v')
endfunc

" The order of the elements are the prioritization of the surrounds
let g:pairList = [['(', ')'], ['[', ']'], ['{', '}'], ['<', '>']]

func! PowerSelection()
  let l:savedPos = getpos('.')
  let l:cachedPair = []

  for l:pair in g:pairList
    if ! s:hasPair(l:pair)
      continue
    endif

    call add(l:cachedPair, l:pair)
    execute "normal! \<esc> vi" . l:pair[0]

    if s:didSelectInline()
      return
    endif
    
    :call setpos('.', l:savedPos)
  endfor

  for [l:left, l:right] in l:cachedPair
    execute "normal! \<esc>f".l:right."F".l:left."vi".l:left

    if s:didSelectInline()
      return
    endif

    :call setpos('.', l:savedPos)
  endfor

  for [l:left, l:right] in l:cachedPair
    execute "normal! \<esc>0f".l:left."vi".l:left

    if s:didSelectInline()
      return
    endif

    :call setpos('.', l:savedPos) 
  endfor
endfunc

