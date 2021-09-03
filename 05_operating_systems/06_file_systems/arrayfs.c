#define FUSE_USE_VERSION 31
#define ROOTDIR_LENGTH 10
#define DIR_LENGTH 10

#include <fuse.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <errno.h>
#include <fcntl.h>
#include <stddef.h>
#include <assert.h>

char *array[ROOTDIR_LENGTH][DIR_LENGTH];

// directory one = {"a", NULL, "b",
// 				 NULL, NULL, NULL,
// 				 NULL, NULL, NULL,
// 				 "v"};
// directory two = {"c", NULL, NULL,
// 				 NULL, NULL, NULL,
// 				 NULL, "d", NULL,
// 				 "e"};
// directory three = {NULL, NULL, NULL,
// 				   "f", "g", "h",
// 				   NULL, NULL, NULL,
// 				   NULL};

// directory array_tmp[ROOTDIR_LENGTH] = {NULL, NULL, NULL,
// 										  one, two, three,
// 										  NULL, NULL, NULL,
// 										  NULL};

static void *arrayfs_init(struct fuse_conn_info *conn,
						  struct fuse_config *cfg)
{
	(void)conn;
	cfg->kernel_cache = 1;
	array[0][0] = "a";
	array[0][1] = NULL;
	array[0][2] = NULL;
	array[0][3] = NULL;
	array[0][4] = NULL;
	array[0][5] = NULL;
	array[0][6] = NULL;
	array[0][7] = NULL;
	array[0][8] = NULL;
	array[0][9] = NULL;
	array[1][0] = NULL;
	array[1][1] = NULL;
	array[1][2] = NULL;
	array[1][3] = NULL;
	array[1][4] = NULL;
	array[1][5] = NULL;
	array[1][6] = NULL;
	array[1][7] = NULL;
	array[1][8] = NULL;
	array[1][9] = NULL;
	array[2][0] = NULL;
	array[2][1] = "b";
	array[2][2] = NULL;
	array[2][3] = NULL;
	array[2][4] = NULL;
	array[2][5] = NULL;
	array[2][6] = NULL;
	array[2][7] = NULL;
	array[2][8] = NULL;
	array[2][9] = NULL;
	array[3][0] = NULL;
	array[3][1] = NULL;
	array[3][2] = NULL;
	array[3][3] = NULL;
	array[3][4] = NULL;
	array[3][5] = NULL;
	array[3][6] = NULL;
	array[3][7] = NULL;
	array[3][8] = NULL;
	array[3][9] = NULL;
	array[4][0] = NULL;
	array[4][1] = NULL;
	array[4][2] = NULL;
	array[4][3] = NULL;
	array[4][4] = NULL;
	array[4][5] = NULL;
	array[4][6] = NULL;
	array[4][7] = "c";
	array[4][8] = NULL;
	array[4][9] = NULL;
	array[5][0] = NULL;
	array[5][1] = NULL;
	array[5][2] = NULL;
	array[5][3] = NULL;
	array[5][4] = NULL;
	array[5][5] = NULL;
	array[5][6] = NULL;
	array[5][7] = NULL;
	array[5][8] = NULL;
	array[5][9] = NULL;
	array[6][0] = NULL;
	array[6][1] = NULL;
	array[6][2] = NULL;
	array[6][3] = NULL;
	array[6][4] = NULL;
	array[6][5] = NULL;
	array[6][6] = NULL;
	array[6][7] = NULL;
	array[6][8] = NULL;
	array[6][9] = NULL;
	array[7][0] = NULL;
	array[7][1] = NULL;
	array[7][2] = NULL;
	array[7][3] = NULL;
	array[7][4] = NULL;
	array[7][5] = NULL;
	array[7][6] = NULL;
	array[7][7] = NULL;
	array[7][8] = NULL;
	array[7][9] = NULL;
	array[8][0] = NULL;
	array[8][1] = "d";
	array[8][2] = NULL;
	array[8][3] = NULL;
	array[8][4] = NULL;
	array[8][5] = NULL;
	array[8][6] = NULL;
	array[8][7] = NULL;
	array[8][8] = NULL;
	array[8][9] = NULL;
	array[9][0] = NULL;
	array[9][1] = NULL;
	array[9][2] = NULL;
	array[9][3] = NULL;
	array[9][4] = NULL;
	array[9][5] = NULL;
	array[9][6] = NULL;
	array[9][7] = NULL;
	array[9][8] = NULL;
	array[9][9] = NULL;

	return array;
}

int parse_dir(const char *path)
{
	if (strlen(path) == 1)
	{
		return -1;
	}
	return path[1] - '0';
}

int parse_file(const char *path)
{
	if (strlen(path) < 4)
	{
		return -1;
	}
	return path[3] - '0';
}

static int arrayfs_getattr(const char *path, struct stat *stbuf,
						   struct fuse_file_info *fi)
{
	(void)fi;
	int res = 0;
	struct fuse_context fc = *fuse_get_context();
	//char *arr[ROOTDIR_LENGTH][DIR_LENGTH] = fc.private_data;

	memset(stbuf, 0, sizeof(struct stat));
	int dir = parse_dir(path);
	int file = parse_file(path);
	if (strcmp(path, "/") == 0)
	{
		stbuf->st_mode = S_IFDIR | 0755;
		stbuf->st_nlink = 2 + ROOTDIR_LENGTH;
	}
	else if (strchr(path + 1, '/') == NULL &&
			 dir != -1 && array[dir] != NULL) // no more slashes in the path,
											  // so it is a directory
	{
		stbuf->st_mode = S_IFDIR | 0755;
		stbuf->st_nlink = 2;
	}
	else if (((int)(strchr(path + 1, '/') - path) == 1) &&
			 (file != -1) &&
			 (array[dir][file] != NULL)) // found second slash, so it is a file
	{
		stbuf->st_mode = S_IFREG | 0444;
		stbuf->st_nlink = 1;
		stbuf->st_size = strlen(array[dir][file]);
	}
	else
		res = -ENOENT;

	return res;
}

static int hello_readdir(const char *path, void *buf, fuse_fill_dir_t filler,
						 off_t offset, struct fuse_file_info *fi,
						 enum fuse_readdir_flags flags)
{
	(void)offset;
	(void)fi;
	(void)flags;

	if (strcmp(path, "/") != 0)
		return -ENOENT;

	filler(buf, ".", NULL, 0, 0);
	filler(buf, "..", NULL, 0, 0);
	filler(buf, "hello", NULL, 0, 0);

	return 0;
}

static int hello_open(const char *path, struct fuse_file_info *fi)
{
	if (strcmp(path + 1, "hello") != 0)
		return -ENOENT;

	if ((fi->flags & O_ACCMODE) != O_RDONLY)
		return -EACCES;

	return 0;
}

static int hello_read(const char *path, char *buf, size_t size, off_t offset,
					  struct fuse_file_info *fi)
{
	size_t len;
	(void)fi;
	if (strcmp(path + 1, "hello") != 0)
		return -ENOENT;

	len = strlen("hello");
	if (offset < len)
	{
		if (offset + size > len)
			size = len - offset;
		memcpy(buf, "hello" + offset, size);
	}
	else
		size = 0;

	return size;
}

static const struct fuse_operations arrayfs_oper = {
	.init = arrayfs_init,
	.getattr = arrayfs_getattr,
	.readdir = hello_readdir,
	.open = hello_open,
	.read = hello_read,
};

int main(int argc, char *argv[])
{
	//struct fuse_args args = FUSE_ARGS_INIT(argc, argv);
	printf("argc! %d\n", argc);
	//printf("argv! %s\n", argv[3])	;
	printf("here!\n");
	printf("here!\n");
	printf("!!!%p\n", &arrayfs_oper);
	printf("there!\n");

	return fuse_main(argc,
					 argv,
					 &arrayfs_oper,
					 NULL);
	//fuse_opt_free_args(&args);
}
