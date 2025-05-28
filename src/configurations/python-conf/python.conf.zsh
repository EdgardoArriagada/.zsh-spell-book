if [[ -d $HOME/.pyenv ]]; then
  ZSB_PYTHON_CONF_DIR=$ZSB_DIR/src/configurations/python-conf

  ${zsb}.prepareLazyLoad \
    $ZSB_PYTHON_CONF_DIR/pyenv.init \
    pyenv python3 poetry pip flask
fi

