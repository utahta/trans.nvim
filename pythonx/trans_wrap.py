from trans import Trans
import vim

g_trans = Trans(vim)


def trans(*args):
    g_trans.trans(*args)


def trans_word(*args):
    g_trans.trans_word(*args)