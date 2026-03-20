return {
	{
		"stevearc/conform.nvim",
		opts = {
			formatters = {
				-- golines = {
				-- 	prepend_args = { "--base-formatter", "gofumpt", "--max-len", "110" },
				-- },
				ruff_format = {
					prepend_args = { "--line-length", "110" },
				},
				xmlstarlet = { "fo" },
			},
			formatters_by_ft = {
				python = { "isort", "ruff_format" },
				lua = { "stylua" },
				go = { "gofmt" },
				sh = { "shfmt" },
				proto = { "buf" },
				c = { "clang-format" },
				cpp = { "clang-format" },
				cc = { "clang-format" },
				xml = { "xmlstarlet" },
			},
		},
	},
}
