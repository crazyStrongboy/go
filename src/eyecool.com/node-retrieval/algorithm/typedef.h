
#ifndef __TYPEDEF_H__
#define __TYPEDEF_H__

#ifndef _WIN32
  #ifndef _BASE_TYPES
	#define _BASE_TYPES

	typedef int BOOL;
	#define TRUE    1
	#define FALSE   0
    typedef unsigned char BYTE;
	typedef BYTE* LPBYTE;
	typedef short SHORT;
    typedef unsigned short WORD;
	typedef int INT;
	typedef unsigned int UINT;
    typedef unsigned int DWORD;
	typedef long LONG;
    typedef unsigned long ULONG;
#ifdef X86_X64
	typedef long INT64;
	typedef unsigned long UINT64;
#else
	typedef long long INT64;
	typedef unsigned long long UINT64;
#endif

	typedef void* HANDLE;
    #define INVALID_HANDLE_VALUE	((void*)(-1))

    typedef struct tagPOINT
    {
        int x, y;
    } POINT;

    typedef struct tagSIZE
    {
        int cx, cy;
    } SIZE;

    typedef struct tagRECT
    {
        int left, top, right, bottom;
    } RECT;
#endif // _BASE_TYPES

#endif

typedef struct FACE_DETECT_RESULT
{
    RECT rcFace;
    POINT ptLeftEye, ptRightEye, ptMouth, ptNose;
    int nAngleYaw, nAnglePitch, nAngleRoll, nQuality;
    char FaceData[512];
} FACE_DETECT_RESULT, *LPFACE_DETECT_RESULT;

typedef struct FACE_PROPERTY_RESULT
{
	unsigned char nGender; /* 性别: 1-male, 0-female */
	unsigned char nAge; /* 年龄 */
	unsigned char nRace; /* 种族: 0-未知, 1-白种人, 2-黄种人, 3-黑种人 */
	unsigned char nSmileLevel; /* 微笑程度: 0-100 */
	unsigned char nBeautyLevel; /* 颜值: 0-100 */
} FACE_PROPERTY_RESULT, *LPFACE_PROPERTY_RESULT;


#endif // __TYPEDEF_H__

